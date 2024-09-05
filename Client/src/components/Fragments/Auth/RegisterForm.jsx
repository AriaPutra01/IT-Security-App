"use client";
import { useRef } from "react";
import Swal from "sweetalert2";
import axios from "axios";
import { Link } from "react-router-dom";

export function RegisterForm() {
  const usernameRef = useRef();
  const emailRef = useRef();
  const passwordRef = useRef();
  // const [errors, setErrors] = useState({});

  const handleRegister = async (e) => {
    e.preventDefault();
    const userData = {
      username: usernameRef.current.value,
      email: emailRef.current.value,
      password: passwordRef.current.value,
      role: document.querySelector('input[name="role"]:checked').value,
    };

    try {
      await axios.post("http://localhost:8080/register", userData);
      Swal.fire({
        icon: "success",
        title: "Tambah User Berhasil",
      }).then(() => {
        window.location.href = "/user";
      });
    } catch (error) {
      Swal.fire({
        icon: "error",
        title: "Tambah User Gagal" + error,
        text: "Silahkan Coba Lagi",
      });
    }
  };

  return (
    <>
      <form onSubmit={handleRegister} className="mx-auto">
        <div className="relative z-0 w-full mb-5 group">
          <input
            ref={usernameRef}
            id="username"
            type="username"
            className="block py-2.5 px-0 w-full text-sm text-gray-900 bg-transparent border-0 border-b-2 border-gray-300 appearance-none  focus:outline-none focus:ring-0 focus:border-blue-600 peer"
            placeholder=" "
            required
          />
          <label
            htmlFor="username"
            className="peer-focus:font-medium absolute text-sm text-gray-500 duration-300 transform -translate-y-6 scale-75 top-3 -z-10 origin-[0] peer-focus:start-0 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto peer-focus:text-blue-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-6"
          >
            Username
          </label>
        </div>
        <div className="relative z-0 w-full mb-5 group">
          <input
            ref={emailRef}
            id="email"
            type="email"
            className="block py-2.5 px-0 w-full text-sm text-gray-900 bg-transparent border-0 border-b-2 border-gray-300 appearance-none  focus:outline-none focus:ring-0 focus:border-blue-600 peer"
            placeholder=" "
            required
          />
          <label
            htmlFor="email"
            className="peer-focus:font-medium absolute text-sm text-gray-500 duration-300 transform -translate-y-6 scale-75 top-3 -z-10 origin-[0] peer-focus:start-0 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto peer-focus:text-blue-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-6"
          >
            Email address
          </label>
        </div>
        <div className="relative z-0 w-full mb-5 group">
          <input
            ref={passwordRef}
            id="password"
            type="password"
            className="block py-2.5 px-0 w-full text-sm text-gray-900 bg-transparent border-0 border-b-2 border-gray-300 appearance-none  focus:outline-none focus:ring-0 focus:border-blue-600 peer"
            placeholder=" "
            required
          />
          <label
            htmlFor="password"
            className="peer-focus:font-medium absolute text-sm text-gray-500 duration-300 transform -translate-y-6 scale-75 top-3 -z-10 origin-[0] peer-focus:start-0 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto peer-focus:text-blue-600 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-6"
          >
            Password
          </label>
        </div>
        <div className="flex gap-10 my-6">
          <div>
            <input
              id="role-user"
              type="radio"
              value="user"
              name="role"
              required
              className="w-4 h-4 text-blue-600 bg-slate-200 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
            />
            <label
              htmlFor="role-user"
              className="ms-2 text-sm font-medium text-gray-500 dark:text-gray-300"
            >
              User
            </label>
          </div>
          <div>
            <input
              id="role-admin"
              type="radio"
              value="admin"
              name="role"
              required
              className="w-4 h-4 text-blue-600 bg-slate-200 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
            />
            <label
              htmlFor="role-admin"
              className="ms-2 text-sm font-medium text-gray-500 dark:text-gray-300"
            >
              Admin
            </label>
          </div>
        </div>  
        <div className="grid grid-cols-5 gap-2 border-t-4 border-slate-600 rounded pt-4">
          <Link
            to="/user"
            type="button"
            className="col-span-2 text-white bg-yellow-700 hover:bg-yellow-800 focus:ring-4 focus:outline-none focus:ring-yellow-300 font-medium rounded-lg text-sm w-full px-5 py-2.5 text-center"
          >
            Kembali
          </Link>
          <button
            type="submit"
            className="col-span-3 text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full px-5 py-2.5 text-center"
          >
            Tambah
          </button>
        </div>
      </form>
    </>
  );
}
