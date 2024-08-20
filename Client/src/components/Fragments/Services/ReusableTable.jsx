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
    <div className="overflow-auto max-w-screen p-2">
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
                      Edit
                    </Button>
                    <Button
                      onClick={() => handleDelete(data.ID)}
                      color="failure"
                    >
                      Delete
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
                  No Data Available
                </Badge>
              </Table.Cell>
            </Table.Row>
          </Table.Body>
        )}
      </Table>
    </div>
  );
};
