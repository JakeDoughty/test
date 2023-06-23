import {
  SET_APPLICATIONS,
  SET_SELECTED_APPLICATION,
  SET_EVENTS,
  SET_SELECTED_EVENT,
  BackendAction,
} from '../actions/backend';
import { Application, Event, List } from '../backend';

export interface BackendState {
  applications: List<Application> | null;
  events: List<Event> | null;
  selectedApplication: Application | null;
  selectedApplicationIndex: number;
  selectedEvent: Event | null;
  selectedEventIndex: number;
}

function emptyState(): BackendState {
  return {
    applications: null,
    events: null,
    selectedApplication: null,
    selectedApplicationIndex: -1,
    selectedEvent: null,
    selectedEventIndex: -1,
  };
}

function setSelectedApplication(state: BackendState, app: Application | null) {
  if (app == null || state.applications === null) {
    state.selectedApplication = null;
    state.selectedApplicationIndex = -1;
  } else {
    state.selectedApplication = app;
    state.selectedApplicationIndex = state.applications.items.findIndex(
      (a) => a.id == app.id
    );
  }
  state.events = null;
  state.selectedEvent = null;
  state.selectedEventIndex = -1;
}

export const backendReducer = (
  state: BackendState,
  action: BackendAction
): BackendState => {
  let newState = { ...state };
  switch (action.type) {
    case SET_APPLICATIONS:
      // changing the application will reset every thing else(selected )
      if (action.value == null) {
        newState = emptyState();
      } else if (state.selectedApplication !== null) {
        // let me see if I have the selected application in new list or not
        const selectedApp = state.selectedApplication;
        action.value.items.findIndex((app) => app.id == selectedApp.id);
      }
      return newState;
    case SET_SELECTED_APPLICATION:
      newState.events = null;
      newState.selectedEvent = null;
      newState.selectedApplicationIndex = -1;
      return {
        ...state,
        selectedApplication: action.value,
      };
    case SET_EVENTS:
      return {
        ...state,
        events: action.value,
      };
    case SET_SELECTED_EVENT:
      return {
        ...state,
        selectedEvent: action.value,
      };
  }
};
