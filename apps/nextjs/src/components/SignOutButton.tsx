"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { goBetterAuthClientBrowser } from "@/lib/gba-client-browser";
import { useQueryClient } from "@tanstack/react-query";

export default function SignOutButton() {
  const queryClient = useQueryClient();

  const router = useRouter();

  const [isLoading, setIsLoading] = useState<boolean>(false);

  const handleSignOut = async () => {
    try {
      setIsLoading(true);
      await goBetterAuthClientBrowser.signOut({});
      queryClient.removeQueries();
      router.push("/auth/sign-in");
    } catch (error: any) {
      setIsLoading(false);
      console.error("Error during sign out:", error);
      toast.error(error.message);
    }
  };

  return (
    <Button
      type="button"
      variant="default"
      className="w-full"
      disabled={isLoading}
      onClick={handleSignOut}
    >
      {isLoading ? "Signing out..." : "Sign Out"}
    </Button>
  );
}
