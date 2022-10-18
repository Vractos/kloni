import axios, { AxiosError, AxiosResponse } from "axios";
import { errorMessages } from '../constants/messages';
import { apiErrorNames } from '../constants/system';

const api = axios.create({
  baseURL: import.meta.env.DEV ? '' : import.meta.env.BASE_API_URL,
  timeout: 60000,
})

api.interceptors.response.use(handleAxiosSuccess, handleAxiosError);

function handleAxiosSuccess(resp: AxiosResponse) {
  return resp;
}

function handleAxiosError(err: AxiosError) {
  console.log(err.request, err.response);
  if (err.response) {
    const { status } = err.response;

    const errorName = translateResponseError(status) as keyof typeof errorMessages;
    throw {
      name: errorName,
      message: errorMessages[errorName],
      status,
    };
  } else if (err.request) {
    const errorName = translateRequestError(err.message) as keyof typeof errorMessages;
    throw {
      name: errorName,
      message: errorMessages[errorName],
    };
  } else {
    throw {
      name: apiErrorNames.UNEXPECTED,
      message: errorMessages.UNEXPECTED,
    };
  }
}
const translateResponseError = (status: number) => {
  switch (status) {
    case 400:
      return apiErrorNames.BAD_REQUEST;
    case 401:
      return apiErrorNames.UNAUTHORIZED;
    case 403:
      return apiErrorNames.FORBIDDEN;
    case 404:
      return apiErrorNames.NOT_FOUND;
    case 409:
      return apiErrorNames.CONFLICT;
    case 422:
      return apiErrorNames.UNPROCESSABLE_ENTITY;
    case 500:
      return apiErrorNames.INTERNAL_SERVER_ERROR;
    case 503:
      return apiErrorNames.SERVICE_UNAVAILABLE;
    default:
      return apiErrorNames.UNKNOWN;
  }
};

const translateRequestError = (requestMessage: string) => {
  if (requestMessage.includes('timeout')) {
    return apiErrorNames.TIMEOUT;
  } else if (requestMessage.includes('network')) {
    return apiErrorNames.NETWORK_CONNECTION;
  } else {
    return apiErrorNames.BAD_REQUEST;
  }
};

export { api };