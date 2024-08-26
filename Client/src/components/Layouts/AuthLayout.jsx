const AuthLayout = (props) => {
  const { header, children } = props;
  return (
    <div className="w-full min-h-screen flex items-center overflow-hidden bg-blue-950">
      <img
        className="flex justify-center items-center z-0 absolute opacity-50 w-full h-full object-cover"
        src="/public/images/freepickBlue.jpg"
        alt="bg"
      />
      <div className="flex w-full justify-center">
        <div className="w-[400px]  p-8 z-10 bg-slate-100 rounded-xl shadow">
          <header>
            <img
              className="mx-auto w-16"
              src="/public/images/logobjb.png"
              alt="bjb"
            />
            <h1 className="text-2xl font-semibold text-center my-4 text-slate-800">
              {header}
            </h1>
          </header>
          <main>{children}</main>
        </div>
      </div>
    </div>
  );
};

export default AuthLayout;
