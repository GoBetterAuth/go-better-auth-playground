import { Navigate, redirect, useNavigate } from "react-router";

import { Avatar, AvatarFallback } from "~/components/ui/avatar";
import { Badge } from "~/components/ui/badge";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { useMe } from "~/hooks/useMe";
import { toast } from "~/hooks/use-toast";
import { Spinner } from "~/components/ui/spinner";

function formatDate(date?: string) {
  if (!date) return "-";
  try {
    return new Date(date).toLocaleString();
  } catch {
    return date;
  }
}

export default function DashboardPage() {
  const navigate = useNavigate();

  const { data, isLoading, isError, error, refetch } = useMe();

  if (isLoading) {
    return (
      <div className="p-6">
        <Spinner />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="p-6">
        <p className="text-red-500">
          Error: {error instanceof Error ? error.message : "Unknown error"}
        </p>
        <Button onClick={() => refetch()}>Retry</Button>
      </div>
    );
  }

  if (!data) {
    return <Navigate to="/auth/sign-in" replace />;
  }

  function signOut() {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    toast({ title: "Signed out", description: "You have been signed out" });
    navigate("/auth/sign-in");
  }

  return (
    <div className="container mx-auto p-6">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-4">
          <Avatar>
            <AvatarFallback>
              {(data?.user.name || data?.user.email || "U")[0].toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div>
            <h1 className="text-2xl font-bold">
              {data
                ? `Welcome, ${data.user.name ?? data.user.email}`
                : "Welcome"}
            </h1>
            <p className="text-sm text-muted-foreground">
              {data?.user.email ?? "—"}
            </p>
          </div>
        </div>

        <div className="flex gap-2">
          <Button
            variant="ghost"
            onClick={() => refetch()}
            disabled={isLoading}
          >
            Refresh
          </Button>
          <Button onClick={signOut}>Sign out</Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader>
            <CardTitle>Email verification</CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <p>Loading…</p>
            ) : data ? (
              <div className="flex items-center gap-2">
                {data.user.emailVerified ? (
                  <Badge variant="secondary">Verified</Badge>
                ) : (
                  <Badge>Unverified</Badge>
                )}
                <span className="text-muted-foreground text-sm">
                  {data.user.email ?? "—"}
                </span>
              </div>
            ) : (
              <p className="text-muted-foreground">No data</p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Member since</CardTitle>
          </CardHeader>
          <CardContent>
            <p>{isLoading ? "…" : formatDate(data?.user.createdAt)}</p>
          </CardContent>
        </Card>
      </div>

      <div className="mt-6">
        <Card>
          <CardHeader>
            <CardTitle>Profile (raw)</CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <p>Loading…</p>
            ) : (
              <pre className="overflow-x-auto text-sm">
                {JSON.stringify(data, null, 2)}
              </pre>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
