/* eslint-disable @typescript-eslint/no-explicit-any */
"use client";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { ButtonLoader } from "@/components/ui/loading-spinner";
import { AlertCircle } from "lucide-react";
import { useAuth } from "@/app/_components/AuthContext";
import { Textarea } from "@/components/ui/textarea";

interface CreateRestaurantFormProps {
  onSubmit: (formData: any) => Promise<void>;
  loading: boolean;
  error: string | null;
}

export function RestaurantForm({
  onSubmit,
  loading,
  error,
}: CreateRestaurantFormProps) {
  const { owner } = useAuth();

  const [formData, setFormData] = useState({
    name: "",
    description: "",
    main_image:
      "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?q=80&w=2070&auto=format&fit=crop",
    images: [],
    location: {
      address: {
        city: "",
        district: "",
      },
      latitude: 0,
      longitude: 0,
    },
    details: {
      capacity: 0,
      opening_hours: "09:00-22:00",
    },
  });

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;

    if (name.includes(".")) {
      const [parent, child] = name.split(".");
      setFormData({
        ...formData,
        [parent]: {
          ...formData[parent as keyof typeof formData],
          [child]: value,
        },
      });
    } else {
      setFormData({
        ...formData,
        [name]: value,
      });
    }
  };

  const handleAddressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    setFormData({
      ...formData,
      location: {
        ...formData.location,
        address: {
          ...formData.location.address,
          [name]: value,
        },
      },
    });
  };

  const handleCoordinateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    setFormData({
      ...formData,
      location: {
        ...formData.location,
        [name]: parseFloat(value) || 0,
      },
    });
  };

  const handleDetailsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    setFormData({
      ...formData,
      details: {
        ...formData.details,
        [name]: name === "capacity" ? parseInt(value) || 0 : value,
      },
    });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (owner) {
      onSubmit(formData);
    }
  };

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>Create New Restaurant</CardTitle>
      </CardHeader>
      <CardContent>
        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Basic Information */}
          <div className="space-y-4">
            <h2 className="text-lg font-medium">Basic Information</h2>

            <div className="space-y-2">
              <Label htmlFor="name">Restaurant Name *</Label>
              <Input
                id="name"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
                placeholder="Enter restaurant name"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description *</Label>
              <Textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleChange}
                required
                placeholder="Describe your restaurant"
                rows={4}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="main_image">Main Image URL *</Label>
              <Input
                id="main_image"
                name="main_image"
                type="url"
                value={formData.main_image}
                onChange={handleChange}
                required
                placeholder="https://example.com/image.jpg"
              />
              <p className="text-sm text-muted-foreground">
                Enter a valid URL for your restaurant's main image
              </p>
            </div>
          </div>

          {/* Location */}
          <div className="space-y-4 pt-4 border-t">
            <h2 className="text-lg font-medium">Location</h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="city">City *</Label>
                <Input
                  id="city"
                  name="city"
                  value={formData.location.address.city}
                  onChange={handleAddressChange}
                  required
                  placeholder="Enter city"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="district">District *</Label>
                <Input
                  id="district"
                  name="district"
                  value={formData.location.address.district}
                  onChange={handleAddressChange}
                  required
                  placeholder="Enter district"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="latitude">Latitude *</Label>
                <Input
                  id="latitude"
                  name="latitude"
                  type="number"
                  step="0.000001"
                  value={formData.location.latitude}
                  onChange={handleCoordinateChange}
                  required
                  placeholder="Enter latitude"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="longitude">Longitude *</Label>
                <Input
                  id="longitude"
                  name="longitude"
                  type="number"
                  step="0.000001"
                  value={formData.location.longitude}
                  onChange={handleCoordinateChange}
                  required
                  placeholder="Enter longitude"
                />
              </div>
            </div>
          </div>

          {/* Details */}
          <div className="space-y-4 pt-4 border-t">
            <h2 className="text-lg font-medium">Restaurant Details</h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="capacity">Capacity (seats) *</Label>
                <Input
                  id="capacity"
                  name="capacity"
                  type="number"
                  min="1"
                  value={formData.details.capacity}
                  onChange={handleDetailsChange}
                  required
                  placeholder="Enter capacity"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="opening_hours">Opening Hours *</Label>
                <Input
                  id="opening_hours"
                  name="opening_hours"
                  value={formData.details.opening_hours}
                  onChange={handleDetailsChange}
                  required
                  placeholder="e.g., 09:00-22:00"
                />
                <p className="text-xs text-muted-foreground">
                  Format: HH:MM-HH:MM
                </p>
              </div>
            </div>
          </div>

          <Button type="submit" className="w-full" disabled={loading} size="lg">
            {loading ? (
              <>
                Creating Restaurant
                <ButtonLoader className="ml-2" />
              </>
            ) : (
              "Create Restaurant"
            )}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
