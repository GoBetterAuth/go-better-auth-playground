import { useState } from "react";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { Spinner } from "~/components/ui/spinner";
import ENV_CONFIG from "~/constants/env-config";
import { goBetterAuthClient } from "~/lib/gba-client";

const formSchema = z.object({
  email: z.email("Invalid email address"),
});

type FormValues = z.infer<typeof formSchema>;

export default function MagicLinkSignInButton() {
  const [isSubmitted, setIsSubmitted] = useState<boolean>(false);
  const [submittedEmail, setSubmittedEmail] = useState<string>("");

  const { register, handleSubmit, reset, formState } = useForm<FormValues>({
    mode: "onChange",
    defaultValues: {
      email: "john.doe@example.com",
    },
    resolver: zodResolver(formSchema),
  });

  const { errors, isSubmitting, isValid } = formState;

  const onSubmit = async (data: FormValues) => {
    try {
      const response = await goBetterAuthClient.magicLink.signIn({
        email: data.email,
        callbackUrl: `${ENV_CONFIG.baseUrl}/auth/magic-link/exchange`,
      });

      setSubmittedEmail(data.email);
      setIsSubmitted(true);
      toast.success(response.message);
    } catch (error: any) {
      console.error("Error sending magic link:", error);
      toast.error(
        error?.message || "Failed to send magic link. Please try again.",
      );
    }
  };

  if (isSubmitted) {
    return (
      <div className="w-full max-w-md mx-auto p-4 text-center">
        <div className="bg-green-50 dark:bg-green-950 border border-green-200 dark:border-green-800 rounded-lg p-6">
          <h3 className="font-semibold text-green-900 dark:text-green-100 mb-2">
            Check your email
          </h3>
          <p className="text-sm text-green-700 dark:text-green-200 mb-4">
            We&apos;ve sent a magic link to <strong>{submittedEmail}</strong>
          </p>
          <Button
            variant="outline"
            onClick={() => {
              setIsSubmitted(false);
              reset();
            }}
            className="mt-2"
          >
            Send to different email
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full">
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="your@email.com"
              {...register("email")}
              aria-invalid={!!errors.email}
            />
            {errors.email ? (
              <em className="text-red-500 text-sm">{errors.email.message}</em>
            ) : null}
          </div>

          <div>
            <Button
              type="submit"
              className="w-full"
              disabled={!isValid || isSubmitting}
            >
              {isSubmitting ? (
                <>
                  <Spinner className="size-4" />
                  Sending link...
                </>
              ) : (
                "Send Magic Link"
              )}
            </Button>
          </div>
        </div>
      </form>
    </div>
  );
}
