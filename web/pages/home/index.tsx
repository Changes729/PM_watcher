import React from "react";

import { useTranslation } from "react-i18next";

export default function Homepage() {
  const homePageDomRef = React.useRef<HTMLDivElement>(null);
  const headerCheckboxRef = React.useRef<HTMLInputElement>(null);
  const { t } = useTranslation();
  const csvColumns = [
    t("device-name"),
    t("equipment-magnification"),
    t("watt"),
    t("adjust-power"),
    t("last-collection"),
  ]; // Specify your CSV headers here
  const testData = [
    {
      name: "Test Device",
      magnification: "10x",
      watt: "100W",
      adjustPower: "5W",
      lastCollection: "2024-06-01 12:00:00",
    },
    {
      name: "Sample Device",
      magnification: "20x",
      watt: "200W",
      adjustPower: "10W",
      lastCollection: "2024-06-01 13:00:00",
    },
    {
      name: "Demo Device",
      magnification: "15x",
      watt: "150W",
      adjustPower: "7W",
      lastCollection: "2024-06-01 14:00:00",
    },
  ];

  function updateData() {

  }

  function handleDownload() {
    const csvString = [
      csvColumns,
      ...testData.map((item) => [
        item.name,
        item.magnification,
        item.watt,
        item.adjustPower,
        item.lastCollection,
      ]),
    ]
      .map((row) => row.join(","))
      .join("\n");

    const currentDate = new Date();
    const formattedDate = `${currentDate.getFullYear()}-${String(
      currentDate.getMonth() + 1
    ).padStart(2, "0")}-${String(currentDate.getDate()).padStart(
      2,
      "0"
    )}_${String(currentDate.getHours()).padStart(2, "0")}-${String(
      currentDate.getMinutes()
    ).padStart(2, "0")}-${String(currentDate.getSeconds()).padStart(2, "0")}`;
    const blob = new Blob([csvString], { type: "text/csv" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `powerMeter-${formattedDate}.csv`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }

  function handleCheckboxChange(event: React.ChangeEvent<HTMLInputElement>) {
    const checkboxes = homePageDomRef.current?.querySelectorAll(
      'input[type="checkbox"]'
    ) as NodeListOf<HTMLInputElement>;
    const isChecked = event.target.checked;

    if (event.target === headerCheckboxRef.current) {
      checkboxes.forEach((checkbox) => {
        checkbox.checked = isChecked;
      });
    } else {
      const allChecked = Array.from(checkboxes).every((checkbox) => {
        if (checkbox === headerCheckboxRef.current) return true;
        else return checkbox.checked;
      });
      headerCheckboxRef.current!.checked = allChecked;
    }
  }

  return (
    <div ref={homePageDomRef}>
      <img src="http://skinspath.acshoes.com/SkinsPath1/201708/7d2688a7-9550-415b-859b-408dd06b853c/Skins/zh-CN/Website/Images/logo.png"></img>
      <h1>{t("Digital-meter-reading-system")} </h1>
      <button onClick={handleDownload}>{t("download")}</button>
      <button>{t("edit-device-info")}</button>
      <table>
        <thead>
          <tr>
            <th scope="col">
              <input
                ref={headerCheckboxRef}
                type="checkbox"
                onChange={handleCheckboxChange}
              />
            </th>
            {csvColumns.map((col, index) => (
              <th key={index} scope="col">
                {col}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {testData.map((device, index) => (
            <tr key={index}>
              <td>
                <input type="checkbox" onChange={handleCheckboxChange} />
              </td>
              <td>{device.name}</td>
              <td>{device.magnification}</td>
              <td>{device.watt}</td>
              <td>{device.adjustPower}</td>
              <td>{device.lastCollection}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
