import { Link, useRouteError } from "react-router-dom";
const ErrorPage = () => {
  const error = useRouteError();
  return (
    <section className="h-screen flex justify-center items-center">
      <div className="py-8 px-4 mx-auto max-w-screen-xl lg:py-16 lg:px-6">
        <div className="mx-auto max-w-screen-sm text-center">
          <h1 className="mb-4 text-6xl tracking-tight font-extrabold text-primary-600 text-red-700">
            {error.statusText || error.message}
          </h1>
          <p className="mb-4 text-3xl tracking-tight font-bold text-slate-100 md:text-4xl ">
            Ada yang hilang.
          </p>
          <p className="mb-4 text-lg font-normal text-slate-400 ">
            Maaf, kami tidak dapat menemukan halaman tersebut. Anda akan
            menemukan banyak hal untuk dijelajahi di halaman utama.
          </p>
          <Link
            to="/"
            className="inline-flex text-slate-100 border-2 border-blue-white hover:text-blue-500 hover:border-blue-500 hover:scale-105 transition-all bg-primary-600 hover:bg-primary-800 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center my-4"
          >
            Kembali ke Halaman Utama
          </Link>
        </div>
      </div>
    </section>
  );
};

export default ErrorPage;
