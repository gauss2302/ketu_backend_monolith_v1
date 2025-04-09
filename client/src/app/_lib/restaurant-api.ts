// src/app/_lib/restaurantApi.ts
import api from "./api";

// Interface for restaurant creation data
export interface CreateRestaurantDTO {
  name: string;
  description: string;
  main_image: string;
  images: string[];
  location: {
    address: {
      city: string;
      district: string;
    };
    latitude: number;
    longitude: number;
  };
  details: {
    capacity: number;
    opening_hours: string;
  };
}

// Interface for restaurant response
export interface RestaurantResponse {
  id: number;
  name: string;
  description: string;
  main_image: string;
  images: string[];
  is_verified: boolean;
  location: {
    address: {
      city: string;
      district: string;
    };
    latitude: number;
    longitude: number;
  };
  details: {
    rating: number;
    capacity: number;
    opening_hours: string;
  };
  created_at: string;
  updated_at: string;
}

// Create a new restaurant with token
export const createRestaurant = async (
  data: CreateRestaurantDTO,
  token: string
): Promise<RestaurantResponse> => {
  try {
    const response = await api.post("/restaurants", data, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    console.error("Error creating restaurant:", error);
    throw error;
  }
};

// Get all restaurants
export const getAllRestaurants = async (
  token: string
): Promise<RestaurantResponse[]> => {
  try {
    const response = await api.get("/restaurants", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data.data;
  } catch (error) {
    console.error("Error fetching restaurants:", error);
    throw error;
  }
};

// Get owner's restaurants
export const getOwnerRestaurants = async (
  token: string
): Promise<RestaurantResponse[]> => {
  try {
    const response = await api.get("/restaurants/my", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data.data;
  } catch (error) {
    console.error("Error fetching owner restaurants:", error);
    throw error;
  }
};

// Get restaurant by ID
export const getRestaurantById = async (
  id: number,
  token: string
): Promise<RestaurantResponse> => {
  try {
    const response = await api.get(`/restaurants/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    console.error(`Error fetching restaurant with ID ${id}:`, error);
    throw error;
  }
};

// Update restaurant
export const updateRestaurant = async (
  id: number,
  data: Partial<CreateRestaurantDTO>,
  token: string
): Promise<RestaurantResponse> => {
  try {
    const response = await api.put(`/restaurants/${id}`, data, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    console.error(`Error updating restaurant with ID ${id}:`, error);
    throw error;
  }
};

// Delete restaurant
export const deleteRestaurant = async (
  id: number,
  token: string
): Promise<void> => {
  try {
    await api.delete(`/restaurants/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  } catch (error) {
    console.error(`Error deleting restaurant with ID ${id}:`, error);
    throw error;
  }
};
