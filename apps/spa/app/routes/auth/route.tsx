import { Navigate, Outlet, useLoaderData } from "react-router";

import { fetchMe } from "~/hooks/useMe";

export async function clientLoader() {
  try {
    const data = await fetchMe();
    return data;
  } catch (error) {
    console.log("showing error");
    console.error(error);
  }

  return null;
}

export default function AuthLayout() {
  const data = useLoaderData<typeof clientLoader>();

  if (data) {
    if (!data.user.email_verified) {
      return <Navigate to="/auth/email-verification" replace />;
    }

    return <Navigate to="/dashboard" replace />;
  }

  return (
    <div className="min-h-screen bg-white">
      <div className="flex items-center justify-center min-h-screen px-4 sm:px-6 lg:px-8">
        <div className="w-full max-w-md">
          {/* Header */}
          <div className="grid place-items-center mb-12 items-center gap-2">
            <img src="/app-logo.png" alt="App Logo" height={100} width={100} />
          </div>

          {/* Auth Content */}
          <div className="bg-slate-50 rounded-lg shadow-lg border border-slate-200 p-8">
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
}
