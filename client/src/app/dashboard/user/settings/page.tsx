"use client";

import React from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useTheme } from "next-themes";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export default function SettingsPage() {
  const { theme, setTheme } = useTheme();
  const [emailNotifications, setEmailNotifications] = React.useState(true);
  const [password, setPassword] = React.useState("");
  const [email, setEmail] = React.useState("");

  const handleEmailNotificationsChange = (checked: boolean) => {
    setEmailNotifications(checked);
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  const handleSaveChanges = () => {
    // In a real application, you would send these changes to your backend.
    console.log("Settings saved:", {
      theme,
      emailNotifications,
      password,
      email,
    });
    alert("Settings saved!");
  };

  return (
    <div className="p-4">
      <Card className="w-full max-w-2xl mx-auto">
        <CardHeader>
          <CardTitle>User Settings</CardTitle>
          <CardDescription>Manage your account settings.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label>Theme</Label>
            <Select value={theme} onValueChange={setTheme}>
              <SelectTrigger className="w-full sm:w-[180px]">
                <SelectValue placeholder="Select theme" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="light">Light</SelectItem>
                <SelectItem value="dark">Dark</SelectItem>
                <SelectItem value="system">System</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between">
            <Label htmlFor="email-notifications">Email Notifications</Label>
            <Switch
              id="email-notifications"
              checked={emailNotifications}
              onCheckedChange={handleEmailNotificationsChange}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Change Password</Label>
            <Input
              type="password"
              id="password"
              value={password}
              onChange={handlePasswordChange}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="email">Change Email</Label>
            <Input
              type="email"
              id="email"
              value={email}
              onChange={handleEmailChange}
            />
          </div>
          <Button onClick={handleSaveChanges}>Save Changes</Button>
        </CardContent>
      </Card>
    </div>
  );
}
