import React, { useEffect, useState } from 'react';

import { Dropdown } from 'primereact/dropdown';

import { Application, BackendConnection, Event } from './backend';
import EventList from './EventList';

interface Props {
  appId: string;
  backend: BackendConnection;
}

function EventsPage({
  backend,
  appId,
  applications,
  onSelectedAppChange,
}: Props) {
  const [pageIndex, setPageIndex] = useState(0);
  const [events, setEvents] = useState<Event[]>([]);
  useEffect(() => {
    backend
      .getEvents(appId, { pageIndex, pageSize: 50 })
      .then((events) => setEvents(events.items));
  });
  const options = 
  return (
    <div>
      <Dropdown
        value={applications}
        placeholder="select an application"
        className="w-full md:w-14rem"
        options={applications}
        onChange={(e) => onSelectedAppChange(e.value)}
      />
      <EventList events={events} />
    </div>
  );
}

export default EventsPage;
