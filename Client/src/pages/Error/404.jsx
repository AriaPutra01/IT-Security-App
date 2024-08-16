import { Link, useRouteError } from "react-router-dom";
const ErrorPage = () => {
  const error = useRouteError();
  return (
    <section className="h-screen flex justify-center items-center bg-white dark:bg-gray-900">
      <div className="py-8 px-4 mx-auto max-w-screen-xl lg:py-16 lg:px-6">
        <div className="mx-auto max-w-screen-sm text-center">
          <h1 className="mb-4 text-8xl tracking-tight font-extrabold text-primary-600 text-red-600 dark:text-primary-500">
            {error.statusText || error.message}
          </h1>
          <p className="mb-4 text-3xl tracking-tight font-bold text-gray-900 md:text-4xl dark:text-white">
            Ada yang hilang.
          </p>
          <p className="mb-4 text-lg font-light text-gray-500 dark:text-gray-400">
            Maaf, kami tidak dapat menemukan halaman itu. Anda akan menemukan
            banyak hal untuk dijelajahi di halaman utama.{" "}
          </p>
          <Link
            to="/"
            className="inline-flex text-black hover:text-blue-600 dark:text-blue-400 dark:hover:text-blue-300 bg-primary-600 hover:bg-primary-800 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:focus:ring-primary-900 my-4"
          >
            Kembali ke Halaman Utama
          </Link>
        </div>
      </div>
    </section>
  );
};

export default ErrorPage;
