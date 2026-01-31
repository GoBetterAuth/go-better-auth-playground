import { Link } from "react-router";

import { ArrowRight, Shield, Key, Mail, UserPlus } from "lucide-react";

export default function IndexPage() {
  const authPages = [
    {
      title: "Sign In",
      description: "Access your account",
      path: "/auth/sign-in",
      icon: Key,
    },
    {
      title: "Sign Up",
      description: "Create a new account",
      path: "/auth/sign-up",
      icon: UserPlus,
    },
    {
      title: "Email Verification",
      description: "Verify your email address",
      path: "/auth/email-verification",
      icon: Mail,
    },
    {
      title: "Reset Password",
      description: "Recover your account",
      path: "/auth/reset-password",
      icon: Shield,
    },
  ];

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4 auth-gradient-bg">
      <div className="w-full max-w-2xl">
        {/* Header */}
        <div className="grid place-items-center mb-12 items-center gap-2">
          <img src="/app-logo.png" alt="App Logo" height={200} width={200} />
        </div>

        {/* Auth Pages Grid */}
        <div className="grid gap-3 sm:grid-cols-2">
          {authPages.map((page) => {
            const Icon = page.icon;
            return (
              <Link
                key={page.path}
                to={page.path}
                className="group bg-card rounded-xl p-5 auth-card border border-border transition-all hover:border-primary/50"
              >
                <div className="flex items-start gap-4">
                  <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center shrink-0">
                    <Icon className="w-5 h-5 text-primary" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <h2 className="font-semibold text-foreground">
                        {page.title}
                      </h2>
                      <ArrowRight className="w-4 h-4 text-muted-foreground opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all" />
                    </div>
                    <p className="text-sm text-muted-foreground mt-1">
                      {page.description}
                    </p>
                  </div>
                </div>
              </Link>
            );
          })}
        </div>
      </div>
    </div>
  );
}
