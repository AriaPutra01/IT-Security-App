import { useState } from "react";

// Hook kustom untuk pagination
export function usePagination(initialPage, itemsPerPage) {
  const [currentPage, setCurrentPage] = useState(initialPage);

  const paginate = (pageNumber) => {
    setCurrentPage(pageNumber);
  };

  const getCurrentPage = () => currentPage;

  const getTotalPages = (totalItems) => Math.ceil(totalItems / itemsPerPage);

  return {
    currentPage,
    paginate,
    getTotalPages,
    getCurrentPage,
  };
}
