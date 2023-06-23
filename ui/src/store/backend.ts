// Here I have used `any` in typescript
// While I go out of my way to remove any from typescript, I think here it is required

export interface Location {
  lat: number;
  long: number;
}
export interface Event {
  type: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data: any;
  createdAt: string;
  applicationID: string;
  sessionID: string;
  location?: Location;
}
export interface Application {
  id: string;
  name: string;
}

export interface PagingInfo {
  pageIndex: number;
  pageSize: number;
}

export interface List<T> {
  items: T[];
  pageIndex?: number;
  pageSize?: number;
  totalItems: number;
  totalPages: number;
}

export type ApplicationList = List<Application>;
export type EventList = List<Event>;

export interface ErrorResponse {
  succeeded: false;
  msg: string;
}
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface SuccessResponse<DataType = any> {
  succeeded: true;
  data: DataType;
}
export type ListReponse<T> = SuccessResponse<List<T>>;

export interface BackendConnection {
  getApplications: (paging: PagingInfo) => Promise<List<Application>>;
  createApplication: (applicationName: string) => Promise<Application>;
  getEvents: (appId: string, paging: PagingInfo) => Promise<List<Event>>;
  getEvent: (eventId: string) => Promise<Event>;
}

export function createBackendConnection(server: string): BackendConnection {
  const apiServer = `${server}/api/v1`;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  async function request<SuccessDataType = any>(
    url: string,
    method = 'GET',
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    body?: any
  ): Promise<SuccessDataType> {
    const options: RequestInit = { method };
    if (body !== undefined) {
      if (typeof body !== 'string') {
        body = JSON.stringify(body);
        options.headers = {
          'Content-Type': 'application/json',
        };
      }
      options.body = body;
    }
    const resp = await fetch(url, options);
    const response = (await resp.json()) as
      | ErrorResponse
      | SuccessResponse<SuccessDataType>;
    if (response.succeeded) {
      return response.data;
    } else {
      throw new Error(response.msg);
    }
  }

  async function getApplications(paging: PagingInfo) {
    const response = await request<List<Application>>(
      `${apiServer}/applications?page=${paging.pageIndex}&pageSize=${paging.pageSize}`
    );
    return response;
  }
  async function createApplication(applicationName: string) {
    const response = await request<Application>(
      `${apiServer}/applications`,
      'POST',
      { name: applicationName }
    );

    return response;
  }
  async function getEvents(appId: string, paging: PagingInfo) {
    const response = await request<List<Event>>(
      `${apiServer}/applications/${appId}/events?pageSize=${paging.pageSize}&page=${paging.pageIndex}`
    );
    return response;
  }
  async function getEvent(eventId: string) {
    const response = await request<Event>(`${apiServer}/events/${eventId}`);
    return response;
  }
  return {
    createApplication,
    getApplications,
    getEvents,
    getEvent,
  };
}
