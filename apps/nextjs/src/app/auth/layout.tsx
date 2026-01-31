"use client";

import { PropsWithChildren } from "react";
import { redirect } from "next/navigation";

import { GetMeResponse } from "go-better-auth";
import { useQuery } from "@tanstack/react-query";

import { Spinner } from "@/components/ui/spinner";
import { goBetterAuthClientBrowser } from "@/lib/gba-client-browser";

export default function AuthLayout({ children }: PropsWithChildren) {
  const { data, isLoading } = useQuery({
    queryKey: ["me"],
    queryFn: async () => {
      try {
        const response = await goBetterAuthClientBrowser.getMe<GetMeResponse>();
        return response;
      } catch (error) {
        console.error(error);
        return null;
      }
    },
  });

  if (isLoading) {
    return (
      <div>
        <Spinner />
      </div>
    );
  }

  if (data) {
    if (!data.user.emailVerified) {
      redirect(`/auth/email-verification?email=${data.user.email}`);
    }
    redirect("/dashboard");
  }

  return <>{children}</>;
}
