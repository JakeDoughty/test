window.track = (function () {
  const defaultTrackServer = 'http://localhost:3000';

  const location = { lat: 0.0, long: 0.0 };
  var positionMonitoringId = null;

  const eventsToReport = [];
  const eventListeners = [];

  var config = null;
  var sessionId = null;
  var lastEventSent = new Date();
  var apiServer = `${defaultTrackServer}/api/v1`;
  var sendingEvent = false;
  var eventReportingInterval = null;

  function isRunning() {
    return eventReportingInterval !== null;
  }

  function xhr(data) {
    if (typeof data === 'string') {
      data = { path: data };
    }
    if (data.method === undefined) {
      data.method = 'GET';
    }

    var x = new XMLHttpRequest();
    x.onreadystatechange = function (ev) {
      if (x.readyState == XMLHttpRequest.DONE) {
        response =
          ['json', 'text', ''].includes(x.responseType) && x.responseText !== ''
            ? JSON.parse(x.responseText)
            : null;
        data.done && data.done(x, ev, response);
      }
    };
    x.onerror = function (e) {
      data.error && data.error(x, e);
    };

    url = `${apiServer}/${data.path}`;
    x.open(data.method, url);

    x.setRequestHeader('accept', 'application/json');
    var body = undefined;
    const requireBody =
      ['PUT', 'POST', 'PATH'].includes(data.method) && data.body !== undefined;
    if (requireBody) {
      body = JSON.stringify(data.body);
      x.setRequestHeader('content-type', 'application/json');
      x.setRequestHeader('browser', navigator.userAgent);
    }
    x.send(body);
  }
  function openSession() {
    if (sessionId === 0 || typeof sessionId === 'string') {
      return; // session is already opened or session is already openning
    }

    // mark that we are openning the session, to avoid re-entry to this function
    sessionId = 0;

    const request = {
      screen: { width: screen.width, height: screen.height },
    };
    if (location.lat != 0) {
      request.location = location;
    }
    xhr({
      method: 'POST',
      path: `applications/${config.applicationID}/sessions`,
      body: request,
      done: function (xhr, e, response) {
        if (xhr.status == 201 || xhr.status == 200) {
          if (response.succeeded) {
            lastEventSent = new Date();
            sessionId = response.data.id;
            //ready && ready();
          } else {
            // this will allow retrying to open a session to the server
            sessionId = false;
            console.log('[ERR] [OpenSession] ' + response.msg);
          }
        }
      },
      error: function (xhr, e) {
        // this will allow retrying to open a session to the server
        sessionId = false;
        console.log(
          '[ERR] [OpenSession] ' +
            (xhr.responseText || 'network operation failed')
        );
      },
    });
  }
  function closeSession() {
    if (typeof sessionId !== 'string') {
      return;
    }

    xhr({
      method: 'DELETE',
      path: `sessions/${sessionId}`,
    });
    sessionId = null;
    eventsToReport.length = 0; // drop all possible pending events
  }
  function sendEventToServer(data) {
    xhr({
      method: 'POST',
      path: `sessions/${sessionId}/events`,
      body: data.event,
      done: function (xhr, e, response) {
        sendingEvent = false;
        lastEventSent = new Date(); // mark the time that we reported next
        data.done && data.done(xhr, e, response);
      },
      error: function (xhr, e) {
        sendingEvent = false;
        data.error && data.error(xhr, e);
      },
    });
  }
  function pushBackEvent(event) {
    if (event.type !== '@ping') {
      eventsToReport.splice(0, 0, event);
    }
  }
  function reportNextEvent() {
    if (typeof sessionId !== 'string') {
      if (!isRunning()) {
        closeSession();
        return;
      }
      openSession();
      return;
    }
    if (sendingEvent) {
      // we are already sending some event
      return;
    }
    if (eventsToReport.length === 0) {
      // nothing to report
      const now = new Date();
      const diffTimeInMS = Math.abs(now - lastEventSent);
      const diffMinutes = diffTimeInMS / (1000 * 60);
      if (diffMinutes > 10) {
        // there is more than 10 minutes that we didn't report any event to the server,
        // to keep the session alive
        eventsToReport.push({ type: 'ping' });
      } else {
        return;
      }
    }

    sendingEvent = true;
    const event = eventsToReport.splice(0, 1)[0];
    if (event.location === undefined && location.lat != 0) {
      event.location = location;
    }

    console.log('reporting event to the server: ', event);
    sendEventToServer({
      event: event,
      done: function (xhr, e, response) {
        if (!isRunning()) {
          return;
        }

        if (xhr.status == 401) {
          sendingEvent = false;
          // our session is not valid any more, try to open a new session with the server
          sessionId = false;
          openSession();
        } else {
          if (xhr.status === 200 || xhr.status === 201) {
            // this event reported successfully
            reportNextEvent();
            lastEventSent = new Date(); // report the time of reporting last event
          } else {
            // failed to report the event, push the event back to the queue and
            // wait for next loop to report it
            pushBackEvent(event);
          }
        }
      },
      error: function (xhr, e) {
        pushBackEvent(event);
        sendingEvent = false;
      },
    });
  }
  function doReportEvent(eventType, eventData) {
    if (isRunning()) {
      const event = { type: eventType, data: eventData || {} };
      if (event.location === undefined && location.lat != 0) {
        event.location = location;
      }
      eventsToReport.push(event);
      reportNextEvent();
    }
  }
  function emptyDataCollector(ev) {
    return {};
  }
  function _addEventListener(target, eventType, reportName, dataCollector) {
    const listener = (ev) => {
      doReportEvent(reportName, dataCollector(ev));
    };

    eventListeners.push({
      type: eventType,
      callback: listener,
      target: target,
    });
    target.addEventListener(eventType, listener);
  }

  return {
    running: isRunning,
    start: function (cfg) {
      if (isRunning()) {
        throw new Error('a session is already started');
      }

      config = cfg;
      apiServer = `${config.server || defaultTrackServer}/api/v1`;
      if (config.withLocation) {
        positionMonitoringId = navigator.geolocation.watchPosition(
          (position) => {
            location.lat = position.coords.latitude;
            location.long = position.coords.longitude;
          }
        );
      }

      // setup event listeners
      // function setupLocationChange() {
      //   window.addEventListener('load', () => {
      //     let oldHref = '';
      //     new MutationObserver((mutations) => {
      //       mutations.forEach(mutation => {

      //       }
      //       if (oldHref !== document.location.href) {
      //         oldHref = document.location.href;
      //         doReportEvent('@page_view', {});
      //       }
      //     }).observe(document.querySelector('body'), {
      //       childList: true,
      //       subtree: true,
      //     });
      //   });
      // }
      const eventsToListen = [
        {
          target: window,
          event: 'popstate',
          name: '@page_view',
          dataCollector: (e) => ({ location: window.location.href }),
        },
      ].concat(config.events || []);
      eventsToListen.forEach((item) => {
        if (typeof item === 'string') {
          item = { target: window, event: item };
        }

        _addEventListener(
          item.target || window,
          item.event,
          item.name || item.type,
          item.dataCollector || emptyDataCollector
        );
      });

      eventReportingInterval = setInterval(reportNextEvent);
      reportNextEvent(); // first tick
    },
    listenToEvents: function (target, event, name, dataCollector) {
      if (isRunning()) {
        _addEventListener(
          target,
          event,
          name || event,
          dataCollector || emptyDataCollector
        );
      }
    },
    repertEvent: function (eventType, eventData) {
      // allow user to directly report their events
      doReportEvent(eventType, eventData);
    },
    stop: function () {
      const n = eventReportingInterval;
      eventReportingInterval = null;
      if (n !== null) {
        clearInterval(n);

        if (positionMonitoringId) {
          navigator.geolocation.clearWatch(positionMonitoringId);
          positionMonitoringId = null;
          location.lat = 0;
          location.long = 0;
        }

        // remove all event listeners
        eventListeners.forEach((item) => {
          item.target.removeEventListener(item.type, item.listener);
        });
        eventListeners.length = 0;

        closeSession();
      }
    },
  };
})();
