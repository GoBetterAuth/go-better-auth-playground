import { useEffect, useRef, useState } from "react";
import { useNavigate, useSearchParams } from "react-router";
import type { JWTTokensResponse } from "go-better-auth/plugins";

import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { Button } from "~/components/ui/button";
import { Spinner } from "~/components/ui/spinner";
import { toast } from "~/hooks/use-toast";
import { goBetterAuthClient } from "~/lib/gba-client";

export default function MagicLinkExchangePage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
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
        const response =
          await goBetterAuthClient.magicLink.exchange<JWTTokensResponse>({
            token: token,
          });

        // Store tokens in localStorage
        localStorage.setItem("accessToken", response.accessToken);
        localStorage.setItem("refreshToken", response.refreshToken);

        setStatus("success");
      } catch (error: any) {
        console.error("Error exchanging magic link token:", error);
        setStatus("error");
        const errorMsg =
          error.message ||
          "Failed to exchange magic link token. Please try again.";
        setErrorMessage(errorMsg);
        toast({
          title: "Exchange Failed",
          description: errorMsg,
        });
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
              <Button onClick={() => navigate("/dashboard")}>
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
                  onClick={() => navigate("/auth/sign-in")}
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
