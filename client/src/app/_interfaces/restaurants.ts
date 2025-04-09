// Restaurant-related interfaces
export interface IAddressDTO {
  city: string;
  district: string;
}

export interface IRestaurantLocationDTO {
  address: IAddressDTO;
  latitude: number;
  longitude: number;
}

export interface IRestaurantDetailsDTO {
  capacity: number;
  opening_hours: string;
}

export interface ICreateRestaurantDTO {
  owner_id: number;
  name: string;
  description: string;
  main_image: string;
  images?: string[];
  location: IRestaurantLocationDTO;
  details: IRestaurantDetailsDTO;
}

export interface ILocationResponse {
  address: IAddressDTO;
  latitude: number;
  longitude: number;
}

export interface IRestaurantDetailsResponse {
  rating: number;
  capacity: number;
  opening_hours: string;
}

export interface IRestaurantResponse {
  id: number;
  name: string;
  description: string;
  main_image: string;
  images: string[];
  is_verified: boolean;
  location: ILocationResponse;
  details: IRestaurantDetailsResponse;
  created_at: string;
  updated_at: string;
}

export interface IRestaurantListResponse {
  data: IRestaurantResponse[];
  pagination: {
    total: number;
    offset: number;
    limit: number;
  };
}
