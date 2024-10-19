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

export interface UnclaimedDonation {
  id: number;
  donorName: string;
  description: string;
}

export interface Donation {
  id: number;
  donor: User;
  recipient: User;
  meal: Meal;
}

export interface DonationClaimSummary {
  claimed: boolean;
  description: string;
  donorName: string;
  recipientName: string;
}

export interface User {
  uuid: string;
  name: string;
}
