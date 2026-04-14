export interface ApiResult<T> {
  statusCode: number;
  error: string;
  data: T;
}
