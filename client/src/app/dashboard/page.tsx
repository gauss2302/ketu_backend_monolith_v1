// // app/dashboard/page.tsx (using withAuth)
// "use client";

// import { Button } from "@/components/ui/button";
// import withAuth from "@/app/_components/withAuth";
// import { useAuth } from "@/app/_components/AuthContext";

// function DashboardPage() {
//   //No need the type
//   const { user, owner, logout, ownerLogout } = useAuth();

//   return (
//     <div className="p-6">
//       <h1 className="text-2xl font-semibold mb-4">Dashboard</h1>
//       {user && (
//         <>
//           <p>Welcome, {user.name}!</p>
//           <p>User ID: {user.id}</p>
//           <p>Email: {user.email}</p>
//           <Button onClick={logout}>Logout</Button>
//         </>
//       )}
//       {owner && (
//         <>
//           <p>Welcome, Owner {owner.name}!</p>
//           <p>Owner ID: {owner.owner_id}</p>
//           <p>Email: {owner.email}</p>
//           <Button onClick={ownerLogout}>Logout</Button>
//         </>
//       )}
//     </div>
//   );
// }

// export default withAuth(DashboardPage);
