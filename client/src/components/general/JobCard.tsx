"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";

interface JobCardProps {
  job: {
    title: string;
    company: {
      name: string;
      logo: string;
    };
    location: string;
    salary: string;
    description: string;
    applyLink?: string;
  };
}

export function JobCard({ job }: JobCardProps) {
  return (
    <Card className="transition-transform duration-200 hover:scale-105">
      {/* Added hover animation */}
      <CardHeader className="flex items-center justify-between space-y-0 pb-2">
        <div className="flex items-center space-x-2">
          <Avatar className="h-8 w-8">
            <AvatarImage src={job.company.logo} alt={job.company.name} />
            <AvatarFallback>{job.company.name.slice(0, 2)}</AvatarFallback>
          </Avatar>
          <div>
            <CardTitle className="text-base">{job.title}</CardTitle>
            <CardDescription className="text-xs">
              {job.company.name}
            </CardDescription>
          </div>
        </div>
        {job.applyLink && (
          <Button variant="link" asChild>
            <a href={job.applyLink} target="_blank" rel="noopener noreferrer">
              Apply
            </a>
          </Button>
        )}
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <p>{job.location}</p>
          <p>{job.salary}</p>
        </div>
        <p className="text-sm mt-2">{job.description}</p>
      </CardContent>
    </Card>
  );
}
