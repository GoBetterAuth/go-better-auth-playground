"use client";

import { useAction } from 'next-safe-action/hooks';

import { Mail } from 'lucide-react';
import { toast } from 'sonner';

import { sendEmailVerificationAction } from '@/app/actions';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function EmailVerificationPage() {
  const { isPending, executeAsync } = useAction(sendEmailVerificationAction);

  const handleSendEmailVerification = async (): Promise<void> => {
    try {
      const data = await executeAsync();
      if (data.serverError) {
        throw new Error(data.serverError);
      }
      if (data.validationErrors) {
        throw new Error(
          Object.entries(data.validationErrors)
            .map(([_, v]) => v)
            .join(", "),
        );
      }
      toast.success("Reset link sent successfully!");
    } catch (error: any) {
      console.error("Error during password reset:", error);
      toast.error(error.message);
    }
  };

  return (
    <div className="h-full w-full p-4 grid place-items-center">
      <Card className="w-full max-w-md mx-auto mt-10">
        <CardHeader className="text-center">
          <Mail className="mx-auto h-12 w-12 text-muted-foreground" />
          <CardTitle>Check Your Email</CardTitle>
        </CardHeader>
        <CardContent className="text-center">
          <p className="text-sm text-muted-foreground mb-4">
            We&#39;ve sent a verification link to your email. Click the link to
            verify your account.
          </p>
          <Button
            variant="outline"
            disabled={isPending}
            onClick={handleSendEmailVerification}
          >
            {isPending ? "Resending..." : "Resend Verification Email"}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
