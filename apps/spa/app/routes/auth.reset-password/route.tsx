import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Mail } from "lucide-react";
import { useNavigate } from "react-router";

import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Field, FieldLabel, FieldError } from "~/components/ui/field";
import { toast } from "~/hooks/use-toast";
import { goBetterAuthClient } from "~/lib/gba-client";

const resetPasswordSchema = z.object({
  email: z.email("Please enter a valid email address"),
});

type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>;

export default function ResetPasswordPage() {
  const navigate = useNavigate();

  const form = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      email: "",
    },
  });

  const onSubmit = async (data: ResetPasswordFormData) => {
    try {
      await goBetterAuthClient.emailPassword.requestPasswordReset({
        email: data.email,
        callbackUrl: "http://localhost:3000/change-password",
      });

      toast({
        title: "Success",
        description: "Password reset link sent.",
      });

      navigate("/auth/sign-in");
    } catch (error: any) {
      toast({
        title: "Reset password failed",
        description: error?.message || "An unknown error occurred",
      });
    }
  };

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
      <Controller
        name="email"
        control={form.control}
        render={({ field, fieldState }) => (
          <Field data-invalid={fieldState.invalid}>
            <FieldLabel htmlFor={field.name}>Email</FieldLabel>
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                {...field}
                id={field.name}
                type="email"
                placeholder="you@example.com"
                className="pl-10"
                aria-invalid={fieldState.invalid}
              />
            </div>
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />

      <Button type="submit" className="w-full mt-6">
        Send reset link
      </Button>
    </form>
  );
}
