import Link from "next/link";

import { GetMeResponse } from "go-better-auth";

import SignOutButton from "@/components/SignOutButton";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { goBetterAuthClientServer } from "@/lib/gba-client-server";

export default async function DashboardPage() {
  let data: GetMeResponse | null = null;
  try {
    data = await goBetterAuthClientServer.getMe();
  } catch (error) {
    console.error(error);
  }
  if (!data) {
    return null;
  }
  const user = data.user;

  return (
    <div className="h-full w-full p-4 grid place-items-center">
      <Card className="max-w-md w-full mx-auto">
        <CardHeader>
          <CardTitle>Welcome, {user.name}!</CardTitle>
          <CardDescription>Your account details</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label className="text-sm font-medium">ID</Label>
            <p className="text-sm text-muted-foreground">{user.id}</p>
          </div>
          <div>
            <Label className="text-sm font-medium">Name</Label>
            <p className="text-sm text-muted-foreground">{user.name}</p>
          </div>
          <div>
            <Label className="text-sm font-medium">Email</Label>
            <p className="text-sm text-muted-foreground">{user.email}</p>
          </div>
          <div>
            <Label className="text-sm font-medium">Email Verified</Label>
            <p className="text-sm text-muted-foreground">
              {user.emailVerified ? "Yes" : "No"}
            </p>
          </div>
          <div>
            <Label className="text-sm font-medium">Created At</Label>
            <p className="text-sm text-muted-foreground">
              {new Date(user.createdAt).toLocaleString()}
            </p>
          </div>
          <div>
            <Label className="text-sm font-medium">Updated At</Label>
            <p className="text-sm text-muted-foreground">
              {new Date(user.updatedAt).toLocaleString()}
            </p>
          </div>
          <div>
            <Label className="text-sm font-medium">Email</Label>
            <Link
              href="/dashboard/email-change"
              className="text-sm underline text-blue-500"
              aria-label="Change email"
            >
              Change email
            </Link>
          </div>
        </CardContent>
        <CardFooter>
          <SignOutButton />
        </CardFooter>
      </Card>
    </div>
  );
}
