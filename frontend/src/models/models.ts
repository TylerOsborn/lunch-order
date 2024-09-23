export interface ApiResult<T> {
    statusCode: number;
    error: string;
    data: T;
}

export interface Meal {
    id: number;
    description: string;
    date: string;
}