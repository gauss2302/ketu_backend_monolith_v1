// app/dashboard/user/page.tsx
"use client";
import { useAuth } from "@/app/_components/AuthContext";
import { JobCard } from "@/components/general/JobCard";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function UserDashboard() {
  const { user, loading, isAuthenticated } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !isAuthenticated()) {
      router.push("/login");
    }
  }, [loading, isAuthenticated, router]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!user) {
    return null;
  }

  // Sample job data
  const jobs = [
    {
      title: "Software Engineer",
      company: {
        name: "Google",
        logo: "https://via.placeholder.com/150",
      },
      location: "Mountain View, CA",
      salary: "$120k - $150k",
      description: "Develop and maintain software applications...",
      applyLink: "https://careers.google.com/jobs/software-engineer",
    },
    {
      title: "Product Manager",
      company: {
        name: "Amazon",
        logo: "https://via.placeholder.com/150",
      },
      location: "Seattle, WA",
      salary: "$130k - $160k",
      description: "Define and drive product strategy...",
      applyLink:
        "https://www.amazon.jobs/en/jobs/2418041/product-manager-technical",
    },
    {
      title: "Product Manager",
      company: {
        name: "Amazon",
        logo: "https://via.placeholder.com/150",
      },
      location: "Seattle, WA",
      salary: "$130k - $160k",
      description: "Define and drive product strategy...",
      applyLink:
        "https://www.amazon.jobs/en/jobs/2418041/product-manager-technical",
    },
    {
      title: "Software Engineer",
      company: {
        name: "Google",
        logo: "https://via.placeholder.com/150",
      },
      location: "Mountain View, CA",
      salary: "$120k - $150k",
      description: "Develop and maintain software applications...",
      applyLink: "https://careers.google.com/jobs/software-engineer",
    },
    {
      title: "Product Manager",
      company: {
        name: "Amazon",
        logo: "https://via.placeholder.com/150",
      },
      location: "Seattle, WA",
      salary: "$130k - $160k",
      description: "Define and drive product strategy...",
      applyLink:
        "https://www.amazon.jobs/en/jobs/2418041/product-manager-technical",
    },
    {
      title: "Product Manager",
      company: {
        name: "Amazon",
        logo: "https://via.placeholder.com/150",
      },
      location: "Seattle, WA",
      salary: "$130k - $160k",
      description: "Define and drive product strategy...",
      applyLink:
        "https://www.amazon.jobs/en/jobs/2418041/product-manager-technical",
    },
  ];

  return (
    <div>
      <h1>User Dashboard</h1>
      <p>Welcome, {user.name}!</p>
      <div className="mt-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {jobs.map((job, index) => (
          <JobCard key={index} job={job} />
        ))}
      </div>
    </div>
  );
}
