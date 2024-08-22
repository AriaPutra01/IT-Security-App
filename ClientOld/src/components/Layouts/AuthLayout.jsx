import { LoginForm, RegisterForm } from "../Fragments/Auth/Form";
const AuthLayout = (props) => {
  const { type } = props;
  return (
    <div className="w-full min-h-screen flex items-center overflow-hidden">
      <img
        className="flex justify-center items-center z-0 absolute w-full h-full object-cover"
        src="/public/images/freepickBlue.jpg"
        alt="bg"
      />
      <div className="w-96 m-auto p-8 z-10 bg-gray-950 bg-opacity-20 rounded-xl shadow">
        <img
          className="mx-auto w-16"
          src="/public/images/logobjbputih.png"
          alt="bjb"
        />
        <h1 className="text-3xl font-bold text-center m-6 text-gray-100">
          {type === "login" ? "Halaman Login" : "Halaman Register"}
        </h1>
        {type === "login" ? <LoginForm /> : <RegisterForm />}
      </div>
    </div>
  );
};

export default AuthLayout;
