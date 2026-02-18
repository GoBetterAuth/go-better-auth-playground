"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { toast } from "sonner";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { goBetterAuthClientBrowser } from "@/lib/gba-client-browser";

export default function MagicLinkExchangePage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get("token");
  const [status, setStatus] = useState<"loading" | "success" | "error">(
    "loading",
  );
  const [errorMessage, setErrorMessage] = useState<string>("");
  const exchangedRef = useRef(false);

  useEffect(() => {
    if (!token || exchangedRef.current) {
      return;
    }

    exchangedRef.current = true;

    const exchangeToken = async () => {
      try {
        await goBetterAuthClientBrowser.magicLink.exchange({
          token: token,
        });
        setStatus("success");
      } catch (error: any) {
        console.error("Error exchanging magic link token:", error);
        setStatus("error");
        setErrorMessage(
          error.message ||
            "Failed to exchange magic link token. Please try again.",
        );
        toast.error(
          error.message ||
            "Failed to exchange magic link token. Please try again.",
        );
      }
    };

    exchangeToken();
  }, [token]);

  return (
    <div className="h-full w-full p-4 grid place-items-center">
      <Card className="w-full max-w-md mx-auto">
        <CardHeader className="text-center">
          {status === "loading" && (
            <>
              <Spinner className="mx-auto size-8 mb-4" />
              <CardTitle>Exchanging Magic Link Token</CardTitle>
            </>
          )}
          {status === "success" && (
            <CardTitle>Successfully Exchanged Magic Link Token!</CardTitle>
          )}
          {status === "error" && (
            <CardTitle className="text-destructive">Exchange Failed</CardTitle>
          )}
        </CardHeader>
        <CardContent className="text-center">
          {status === "loading" && (
            <p className="text-sm text-muted-foreground">
              Please wait while we exchange your magic link token...
            </p>
          )}
          {status === "success" && (
            <>
              <p className="text-sm text-muted-foreground mb-4">
                You have successfully exchanged your magic link token.
              </p>
              <Button onClick={() => router.replace("/dashboard")}>
                Go to Dashboard
              </Button>
            </>
          )}
          {status === "error" && (
            <>
              <p className="text-sm text-destructive mb-4">{errorMessage}</p>
              <div className="flex gap-2 justify-center">
                <Button
                  variant="outline"
                  onClick={() => router.replace("/auth/sign-in")}
                >
                  Back to Sign In
                </Button>
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
