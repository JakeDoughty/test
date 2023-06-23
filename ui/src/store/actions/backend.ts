import { Event, Application, EventList, ApplicationList } from '../backend';

export const SET_APPLICATIONS = 'SET_APPLICATIONS';
export const SET_SELECTED_APPLICATION = 'SET_SELECTED_APPLICATION';
export const SET_EVENTS = 'SET_EVENTS';
export const SET_SELECTED_EVENT = 'SET_SELECTED_EVENT';

export interface SET_APPLICATIONS_Action {
  type: 'SET_APPLICATIONS';
  value: ApplicationList | null;
}
export interface SET_SELECTED_APPLICATION_Action {
  type: 'SET_SELECTED_APPLICATION';
  value: Application | null;
}
export interface SET_EVENTS_Action {
  type: 'SET_EVENTS';
  value: EventList | null;
}
export interface SET_SELECTED_EVENT_Action {
  type: 'SET_SELECTED_EVENT';
  value: Event | null;
}

export type BackendAction =
  | SET_APPLICATIONS_Action
  | SET_SELECTED_APPLICATION_Action
  | SET_EVENTS_Action
  | SET_SELECTED_EVENT_Action;

export const setApplications = (
  applications: ApplicationList | null = null
) => ({
  type: SET_APPLICATIONS,
  value: applications,
});
export const setSelectedApplication = (
  application: Application | null = null
) => ({
  type: SET_SELECTED_APPLICATION,
  value: application,
});
export const setEvents = (events: EventList | null = null) => ({
  type: SET_EVENTS,
  value: events,
});
export const setSelectedEvent = (event: Event | null = null) => ({
  type: SET_SELECTED_EVENT,
  value: event,
});
