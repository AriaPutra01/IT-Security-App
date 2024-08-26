import { LoginForm } from "../../components/Fragments/Auth/LoginForm";
import AuthLayout from "../../components/Layouts/AuthLayout";

export const LoginPage = () => {
  return (
    <AuthLayout header="Halaman Login">
      <LoginForm />
    </AuthLayout>
  );
};
