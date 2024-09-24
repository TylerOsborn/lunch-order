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

export interface Donation {
    id: number;
    name: string;
    description: string;
    date: string;
}