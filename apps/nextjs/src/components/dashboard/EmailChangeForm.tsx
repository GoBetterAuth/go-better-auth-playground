"use client";

import { useState } from 'react';

import { useAction } from 'next-safe-action/hooks';
import Link from 'next/link';

import { useForm } from '@tanstack/react-form';
import { CheckCircle2, Mail } from 'lucide-react';
import { toast } from 'sonner';
import { z } from 'zod';

import { emailChangeAction } from '@/app/actions';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

const formSchema = z
  .object({
    newEmail: z.email("Invalid email address").trim(),
    confirmEmail: z.email("Invalid email address").trim(),
  })
  .refine((data) => data.newEmail === data.confirmEmail, {
    message: "Emails don't match",
    path: ["confirmEmail"],
  });

export default function EmailChangeForm() {
  const [verificationSent, setVerificationSent] = useState(false);
  const [sentEmail, setSentEmail] = useState("");

  const { executeAsync } = useAction(emailChangeAction);

  const form = useForm({
    defaultValues: {
      newEmail: "",
      confirmEmail: "",
    },
    validators: {
      onChange: formSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        const data = await executeAsync(value);
        if (data.serverError) {
          throw new Error(data.serverError);
        }
        if (data.validationErrors) {
          throw new Error(
            Object.entries(data.validationErrors)
              .map(([, v]) => v)
              .join(", "),
          );
        }
        setSentEmail(value.newEmail);
        setVerificationSent(true);
      } catch (error: unknown) {
        const message =
          error instanceof Error ? error.message : "Unknown error occurred";
        toast.error(message);
      }
    },
  });

  if (verificationSent) {
    return (
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="flex justify-center mb-4">
            <CheckCircle2 className="w-12 h-12 text-green-500" />
          </div>
          <CardTitle className="text-2xl font-bold">
            Verification Email Sent
          </CardTitle>
          <CardDescription className="mt-2">
            A verification email has been sent to <strong>{sentEmail}</strong>.
            Please check your inbox and follow the link to confirm your email
            change.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-start gap-3 rounded-lg bg-blue-50 p-3">
              <Mail className="w-5 h-5 text-blue-600 mt-0.5 shrink-0" />
              <div className="text-sm text-blue-800">
                <p className="font-medium">Email not received?</p>
                <p className="text-xs mt-1">
                  Check your spam or junk folder. The link will expire soon.
                </p>
              </div>
            </div>
            <div className="text-center text-sm">
              <Link href="/dashboard" className="text-blue-600 underline">
                Return to Dashboard
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="w-full max-w-md">
      <CardHeader className="text-center">
        <CardTitle className="text-2xl font-bold">
          Request Email Change
        </CardTitle>
        <CardDescription>
          Enter your new email address to receive an email to confirm the change
          to your account.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
          className="grid gap-4"
        >
          <form.Field
            name="newEmail"
            children={(field) => (
              <div className="grid gap-2">
                <Label htmlFor={field.name}>New Email</Label>
                <Input
                  id={field.name}
                  type="email"
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

          <form.Field
            name="confirmEmail"
            children={(field) => (
              <div className="grid gap-2">
                <Label htmlFor={field.name}>Confirm Email</Label>
                <Input
                  id={field.name}
                  type="email"
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
              <Button type="submit" className="w-full" disabled={!canSubmit}>
                {isSubmitting ? "Confirming..." : "Confirm"}
              </Button>
            )}
          />

          <div className="mt-4 text-center">
            <Link href="/dashboard" className="underline text-sm">
              Back to Dashboard
            </Link>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
