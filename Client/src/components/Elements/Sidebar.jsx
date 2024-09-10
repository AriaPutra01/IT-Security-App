import { MoreVertical, ChevronLast, ChevronFirst } from "lucide-react";
import { useContext, createContext, useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom"; // Tambahkan useLocation
import {
  MdKeyboardArrowUp,
  MdKeyboardArrowDown,
  MdKeyboardArrowRight,
} from "react-icons/md";

const SidebarContext = createContext();

export default function Sidebar({ children, img, title, username, email }) {
  const [expanded, setExpanded] = useState(false);

  return (
    <aside className="h-screen">
      <nav className="h-full flex flex-col bg-white border-r shadow-sm transition-all ">
        <div className="py-3 px-4 flex justify-between items-center">
          <div className="flex items-center gap-x-2">
            <img
              src={img}
              className={`overflow-hidden transition-all ${
                expanded ? "w-8" : "w-0"
              }`}
              alt=""
            />
            {expanded && (
              <span className="text-gray-600 font-bold text-xl">{title}</span>
            )}
          </div>
          <button
            onClick={() => {
              setExpanded((curr) => !curr);
            }}
            className="p-1.5 rounded-lg bg-sky-50 hover:bg-sky-100 transition-colors"
          >
            {expanded ? <ChevronFirst /> : <ChevronLast />}
          </button>
        </div>

        <SidebarContext.Provider value={{ expanded }}>
          <ul className="flex-1 px-3">{children}</ul>
        </SidebarContext.Provider>

        <div
          className={`border-t-2 flex justify-${
            expanded ? "start" : "center"
          } items-center p-3`}
        >
          <img
            src={img}
            alt=""
            className={`${expanded ? "w-0" : "w-8"} transition-all rounded-md`}
          />

          <div
            className={`
              flex justify-between items-center
              overflow-hidden transition-all ${expanded ? "w-44 ml-4" : "w-0"}
          `}
          >
            <div className="leading-4">
              <h4 className="font-semibold">{username}</h4>
              <span className="text-xs text-gray-600">{email}</span>
            </div>
            {/* <MoreVertical size={20} /> */}
          </div>
        </div>
      </nav>
    </aside>
  );
}

export function SidebarItem({
  onClick,
  href,
  icon = <MdKeyboardArrowRight />,
  text,
  alert,
}) {
  const { expanded } = useContext(SidebarContext);
  const location = useLocation(); // Dapatkan path saat ini
  const active = location.pathname === href; // Tentukan apakah item aktif

  return (
    <li>
      <Link
        onClick={onClick}
        className={`
        relative flex items-center py-2 px-3 my-1
        font-medium rounded-md cursor-pointer
        transition-colors group
        ${
          active
            ? "bg-gradient-to-tr from-sky-600 to-sky-500 border-b-2 text-slate-100"
            : "hover:bg-sky-50 text-gray-600"
        }
    `}
        to={href}
      >
        {icon}
        <span
          className={`overflow-hidden transition-all ${
            expanded ? "w-44 ml-4" : "w-0"
          }`}
        >
          {text}
        </span>

        {alert && (
          <div
            className={`absolute right-2 w-2 h-2 rounded bg-sky-400 ${
              expanded ? "" : "top-2"
            }`}
          />
        )}

        {!expanded && (
          <div
            className={`
          absolute z-10 left-full rounded-md px-2 py-1 ml-6
          bg-sky-100 text-sky-800 hover:scale-110 hover:ring-2 text-sm
          invisible opacity-20 -translate-x-3 transition-all
          group-hover:visible group-hover:opacity-100 group-hover:translate-x-0
          `}
          >
            {text}
          </div>
        )}
      </Link>
    </li>
  );
}

export function SidebarCollapse({ children, icon, text, alert }) {
  const { expanded } = useContext(SidebarContext);
  const [drop, setDrop] = useState(false);
  const location = useLocation(); // Dapatkan path saat ini

  // Periksa apakah salah satu children aktif
  const isChildActive = children.some(
    (child) => location.pathname === child.props.href
  );

  useEffect(() => {
    if (isChildActive) {
      setDrop(true);
    }
  }, [isChildActive]);

  return (
    <>
      <li
        onClick={() => {
          if (expanded) {
            setDrop((curr) => !curr);
          }
        }}
        className={`relative flex items-center py-2 px-3 my-1 font-medium rounded-md cursor-pointer transition-colors group ${
          isChildActive
            ? "bg-gradient-to-tr from-sky-200 to-sky-100 text-sky-800"
            : "hover:bg-sky-50 text-gray-600"
        }`}
      >
        {icon}

        <span
          className={`overflow-hidden transition-all ${
            expanded ? "w-44 ml-4" : "w-0"
          }`}
        >
          {text}
        </span>

        {expanded && (drop ? <MdKeyboardArrowUp /> : <MdKeyboardArrowDown />)}

        {alert && (
          <div
            className={`absolute right-2 w-2 h-2 rounded bg-sky-400 ${
              expanded ? "" : "top-2"
            }`}
          />
        )}

        {!expanded &&
          children.map((child, index) => {
            return (
              <div key={index} className="relative group">
                <Link
                  to={child.props.href}
                  className={`absolute z-10 left-full rounded-md px-2 py-1 ml-6 bg-sky-100 text-sky-800 hover:scale-110 hover:ring-2 text-sm invisible opacity-20 -translate-x-3 transition-all group-hover:visible group-hover:opacity-100 group-hover:translate-x-0`}
                  style={{ marginTop: `${index * 2}rem` }}
                >
                  <span className="w-max">
                    {child.props.text.replace(/\s+/g, "")}
                  </span>
                </Link>
              </div>
            );
          })}
      </li>

      {expanded && (
        <ul
          className={`transition-max-height ease-in-out overflow-hidden ${
            drop ? "max-h-screen opacity-100" : "max-h-0 opacity-0"
          } ms-4 space-y-2`}
        >
          {children}
        </ul>
      )}
    </>
  );
}
