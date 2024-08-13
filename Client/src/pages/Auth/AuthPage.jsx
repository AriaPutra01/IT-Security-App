import AuthLayout from "../../components/Layouts/AuthLayout";
const LoginPage = () => {
  return <AuthLayout type="login" />;
};

const RegisterPage = () => {
  return <AuthLayout type="register" />;
};

export { LoginPage, RegisterPage };
