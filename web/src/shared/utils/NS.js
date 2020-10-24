export default class NS {
  /**
   * Create a point.
   * @param {'INIT'|'LOADING'|'SUCCESS'|'ERROR'} status
   * @param {string|null} message
   */
  constructor(
    status,
    message,
    statusCode = 0,
    responseTime = 0,
    requestId,
    cached = false,
    hasData = false,
    errorCaught = false
  ) {
    this.code = status;
    this.message = message;
    this.statusCode = statusCode;
    this.hasData = hasData;
    this.responseTime = responseTime;
    this.requestId = requestId;
    this.cached = cached;
    this.errorCaught = errorCaught;
  }

  get isInit() {
    return this.code === "INIT";
  }

  get isLoading() {
    return this.code === "LOADING";
  }

  get isError() {
    this.errorCaught = true;
    return this.code === "ERROR";
  }

  get isSuccess() {
    return this.code === "SUCCESS";
  }

  clone(
    status,
    message,
    statusCode,
    responseTime,
    requestId,
    cached,
    hasData,
    errorCaught
  ) {
    return new NS(
      status === undefined ? this.code : status,
      message === undefined ? this.message : message,
      statusCode === undefined ? this.statusCode : statusCode,
      responseTime === undefined ? this.responseTime : responseTime,
      requestId === undefined ? this.requestId : requestId,
      cached === undefined ? this.cached : cached,
      hasData === undefined ? this.hasData : hasData,
      errorCaught === undefined ? this.errorCaught : errorCaught
    );
  }

  toString() {
    return this.code;
  }
}
