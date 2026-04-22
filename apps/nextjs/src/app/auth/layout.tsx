"use client";

import { PropsWithChildren } from "react";
import { redirect, usePathname } from "next/navigation";

import { useQuery } from "@tanstack/react-query";
import { GetMeResponse } from "authula";

import { Spinner } from "@/components/ui/spinner";
import { authulaClientBrowser } from "@/lib/authula-client-browser";

export default function AuthLayout({ children }: PropsWithChildren) {
  const pathname = usePathname();

  const { data, isLoading } = useQuery({
    queryKey: ["me"],
    queryFn: async () => {
      try {
        const response = await authulaClientBrowser.getMe<GetMeResponse>();
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
      if (pathname === "/auth/email-verification") {
        return <>{children}</>;
      }

      redirect(`/auth/email-verification?email=${data.user.email}`);
    }
    redirect("/dashboard");
  }

  return <>{children}</>;
}
