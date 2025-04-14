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
