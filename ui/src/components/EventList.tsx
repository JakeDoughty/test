import React from 'react';
import { DataTable } from 'primereact/datatable';
import { Column } from 'primereact/column';

import { Event } from './backend';

interface Props {
  events: Event[];
}

const EventList = ({ events }: Props) => {
  return (
    <DataTable value={events}>
      <Column field="type" header="Type"></Column>
      <Column field="createdAt" header="Time"></Column>
    </DataTable>
  );
};

export default EventList;
