import React from "react";
import { Mail } from "lucide-react";

import { Button } from "~/components/ui/button";
import { toast } from "~/hooks/use-toast";
import { goBetterAuthClient } from "~/lib/gba-client";

export default function EmailVerificationPage() {
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await goBetterAuthClient.emailPassword.sendEmailVerification({
        email: localStorage.getItem("email") ?? "",
        callbackUrl: "http://localhost:3000/dashboard",
      });

      toast({
        title: "Verification email sent",
        description:
          "A verification email has been sent to your email address.",
      });
    } catch (error: any) {
      toast({
        title: "Email verification failed",
        description: error?.message || "An unknown error occurred",
      });
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="flex items-start gap-3">
        <Mail className="h-5 w-5 text-muted-foreground" />
        <div>
          <p className="text-sm text-muted-foreground">
            A verification link was sent to your email address.
          </p>
        </div>
      </div>

      <Button type="submit" className="w-full mt-6">
        Send verification email
      </Button>
    </form>
  );
}
