import { Badge, Button, Checkbox, Table } from "flowbite-react";
import { FormatDate } from "../../../Utilities/FormatDate";
import React from "react";

export const ReusableTable = (props) => {
  const {
    formConfig,
    Paginated,
    handleEdit,
    handleDelete,
    selectedIds,
    handleCheckboxChange,
  } = props;

  const renderCellContent = (field, value) => {
    switch (field.type) {
      case "number":
        return `Rp. ${new Intl.NumberFormat("id-ID").format(value)}`;
      case "date":
        return FormatDate(value);
      default:
        return value;
    }
  };

  return (
    <div className="overflow-auto w-full rounded-lg p-2">
      <div className="w-max min-w-full rounded-lg">
        <Table hoverable>
          <Table.Head>
            <Table.HeadCell className="text-center">
              <span>Select</span>
            </Table.HeadCell>
            {formConfig.fields.map((field, index) => (
              <Table.HeadCell key={index}>{field.label}</Table.HeadCell>
            ))}
            <Table.HeadCell>
              <span>Action</span>
            </Table.HeadCell>
          </Table.Head>
          {Paginated.length > 0 ? (
            <Table.Body className="divide-y">
              {Paginated.map((data) => (
                <Table.Row key={data.ID}>
                  <Table.Cell className="text-center">
                    <Checkbox
                      checked={selectedIds.includes(data.ID)}
                      onChange={() => handleCheckboxChange(data.ID)}
                    />
                  </Table.Cell>
                  {formConfig.fields.map((field, index) => (
                    <Table.Cell key={index}>
                      {renderCellContent(field, data[field.name])}
                    </Table.Cell>
                  ))}
                  <Table.Cell>
                    <div className="flex gap-2">
                      <Button
                        onClick={() => handleEdit(data)}
                        action="edit"
                        color="warning"
                      >
                        <svg
                          className="w-6 h-6"
                          aria-hidden="true"
                          xmlns="http://www.w3.org/2000/svg"
                          width="24"
                          height="24"
                          fill="none"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke="currentColor"
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="m14.304 4.844 2.852 2.852M7 7H4a1 1 0 0 0-1 1v10a1 1 0 0 0 1 1h11a1 1 0 0 0 1-1v-4.5m2.409-9.91a2.017 2.017 0 0 1 0 2.853l-6.844 6.844L8 14l.713-3.565 6.844-6.844a2.015 2.015 0 0 1 2.852 0Z"
                          />
                        </svg>
                      </Button>
                      <Button
                        onClick={() => handleDelete(data.ID)}
                        color="failure"
                      >
                        <svg
                          className="w-6 h-6"
                          aria-hidden="true"
                          xmlns="http://www.w3.org/2000/svg"
                          width="24"
                          height="24"
                          fill="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            fillRule="evenodd"
                            d="M8.586 2.586A2 2 0 0 1 10 2h4a2 2 0 0 1 2 2v2h3a1 1 0 1 1 0 2v12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V8a1 1 0 0 1 0-2h3V4a2 2 0 0 1 .586-1.414ZM10 6h4V4h-4v2Zm1 4a1 1 0 1 0-2 0v8a1 1 0 1 0 2 0v-8Zm4 0a1 1 0 1 0-2 0v8a1 1 0 1 0 2 0v-8Z"
                            clipRule="evenodd"
                          />
                        </svg>
                      </Button>
                    </div>
                  </Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
          ) : (
            <Table.Body className="divide-y">
              <Table.Row>
                <Table.Cell
                  colSpan={formConfig.fields.length + 2}
                  className="text-center"
                >
                  <Badge className="p-4 font-bold" color="red">
                    Tidak ada Data
                  </Badge>
                </Table.Cell>
              </Table.Row>
            </Table.Body>
          )}
        </Table>
      </div>
    </div>
  );
};
