"use client";

import { useState } from "react";

import { useForm } from "@tanstack/react-form";
import { toast } from "sonner";
import { z } from "zod";

import { ENV_CONFIG } from "@/constants/env-config";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Spinner } from "@/components/ui/spinner";
import { goBetterAuthClientBrowser } from "@/lib/gba-client-browser";

const formSchema = z.object({
  email: z.email("Invalid email address"),
});

export default function MagicLinkSignInButton() {
  const [isSubmitted, setIsSubmitted] = useState<boolean>(false);
  const [submittedEmail, setSubmittedEmail] = useState<string>("");

  const form = useForm({
    defaultValues: {
      email: "john.doe@example.com",
    },
    validators: {
      onChange: formSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        const response = await goBetterAuthClientBrowser.magicLink.signIn({
          email: value.email,
          callbackUrl: `${ENV_CONFIG.baseUrl}/auth/magic-link/exchange`,
        });
        setSubmittedEmail(value.email);
        setIsSubmitted(true);
        toast.success(response.message);
      } catch (error: any) {
        console.error("Error sending magic link:", error);
        toast.error(
          error.message || "Failed to send magic link. Please try again.",
        );
      }
    },
  });

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
              form.reset();
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
      <form
        onSubmit={(e) => {
          e.preventDefault();
          e.stopPropagation();
          form.handleSubmit();
        }}
      >
        <div className="flex flex-col gap-4">
          <form.Field
            name="email"
            children={(field) => (
              <div className="flex flex-col gap-2">
                <Label htmlFor={field.name}>Email</Label>
                <Input
                  id={field.name}
                  type="email"
                  placeholder="your@email.com"
                  value={field.state.value}
                  onBlur={field.handleBlur}
                  onChange={(e) => field.handleChange(e.target.value)}
                />
                {field.state.meta.isTouched && !field.state.meta.isValid ? (
                  <em className="text-red-500 text-sm">
                    {field.state.meta.errors
                      .map((error) => error?.message)
                      .join(", ")}
                  </em>
                ) : null}
              </div>
            )}
          />

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
            children={([canSubmit, isSubmitting]) => (
              <Button
                type="submit"
                className="w-full"
                disabled={!canSubmit || isSubmitting}
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
            )}
          />
        </div>
      </form>
    </div>
  );
}
