import React, { Fragment } from "react";
import { Badge, Table, Label } from "flowbite-react";
import App from "../../../components/Layouts/App";

export const BaseProjectPage = () => {
  return (
    <App services="Base Project">
      <div className="flex flex-col gap-2 m-2">
        <Label
          value="Project Code Divisi IT Security"
          className="flex justify-center font-black  "
        />
        <Label value="Nomor urut /YYYYY/ZZZZZ/AAAAA/B/YEARS" className="" />
        <Label value="Z" />
        <Table hoverable>
          <Table.Head>
            <Table.HeadCell>Infrastruktur Type</Table.HeadCell>
            <Table.HeadCell>Code</Table.HeadCell>
          </Table.Head>
          <Table.Body>
            <Table.Row>
              <Table.Cell>Software</Table.Cell>
              <Table.Cell>SOF</Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>Hardware</Table.Cell>
              <Table.Cell>HAR</Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>Jasa/Human Resource</Table.Cell>
              <Table.Cell>SER</Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>
        <Label value="A" />
        <Table hoverable>
          <Table.Head>
            <Table.HeadCell>Type Anggaran</Table.HeadCell>
            <Table.HeadCell>Code</Table.HeadCell>
          </Table.Head>
          <Table.Body>
            <Table.Row>
              <Table.Cell>RBB</Table.Cell>
              <Table.Cell>RBB</Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>NON-RBB</Table.Cell>
              <Table.Cell>NRBB</Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>
        <Label value="B" />
        <Table hoverable>
          <Table.Head>
            <Table.HeadCell>NEW PRODUCT</Table.HeadCell>
            <Table.HeadCell>A</Table.HeadCell>
          </Table.Head>
          <Table.Body>
            <Table.Row>
              <Table.Cell>RENEWAL</Table.Cell>
              <Table.Cell>B</Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>
        <Label value="Y" />
        <Table hoverable>
          <Table.Head>
            <Table.HeadCell>GROUP</Table.HeadCell>
            <Table.HeadCell>Code</Table.HeadCell>
          </Table.Head>
          <Table.Body>
            <Table.Row>
              <Table.Cell>Security Operation</Table.Cell>
              <Table.Cell>ITS-ISO</Table.Cell>
            </Table.Row>
            <Table.Row>
              <Table.Cell>Security Architecture & Governence</Table.Cell>
              <Table.Cell>ITS-SAG</Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>
        <Badge className="w-fit p-2">
          <Label value="EXAMPLE :" />
        </Badge>
        <Badge className="w-fit p-2">
          <Label value="0001/ITS-SAG/SOF/RBB/A/2024" />
        </Badge>
      </div>
    </App>
  );
};
