export interface ApiResult<T> {
    statusCode: number;
    error: string;
    data: T;
}


export interface MealType {
    id: number;
    description: string;
}

export interface Meal {
    id: number;
    typeId: number;
    date: string;
}

export interface MenuItem {
    date: string;
    description: string;
}