import { Navigate } from "react-router-dom";
import { LoginForm } from "../../components/Fragments/Auth/LoginForm";
import { useToken } from "../../context/TokenContext";
import AuthLayout from "../../components/Layouts/AuthLayout";

export const LoginPage = () => {
  const { token } = useToken();
  if (token) return <Navigate to="/dashboard" />;
  return (
    <AuthLayout header="Halaman Login">
      <LoginForm />
    </AuthLayout>
  );
};
