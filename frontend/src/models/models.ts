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
  donorName: string;
  description: string;
}

export interface DonationClaimSummary {
  claimed: boolean;
  description: string;
  donorName: string;
  recipientName: string;
}

export interface DonationRequest {
  id: number;
  requesterName: string;
  mealIds: number[];
  status: 'pending' | 'fulfilled' | 'cancelled';
}

export interface MealPreference {
  id: number;
  description: string;
  selected: boolean;
}

export interface MealOrder {
  id: number;
  weekStartDate: string;
  mondayMeal: Meal | null;
  tuesdayMeal: Meal | null;
  wednesdayMeal: Meal | null;
  thursdayMeal: Meal | null;
}

export interface MealOrderRequest {
  weekStartDate: string;
  mondayMealId: number | null;
  tuesdayMealId: number | null;
  wednesdayMealId: number | null;
  thursdayMealId: number | null;
}
