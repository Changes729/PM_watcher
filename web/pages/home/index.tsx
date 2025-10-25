import React, { useEffect } from "react";

import { useTranslation } from "react-i18next";

type DeviceData = {
  ID: string;
  Name: string;
  Magnification: number;
  Watt: number;
  LatestUpdate: string;
};

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
  ];
  var [csvDatas, setCsvDatas] = React.useState([] as DeviceData[]);

  useEffect(() => {
    updateData();
    const interval = setInterval(() => {
      updateData();
    }, 20000); // Update every 5 seconds

    return () => clearInterval(interval); // Cleanup on unmount
  }, []);

  function updateData() {
    fetch("/api/power", {
      method: "GET",
    }).then((res) => {
      res.json().then((data) => {
        console.log(data);
        setCsvDatas(data);
      });
    });
  }

  function handleDownload() {
    const csvString = [
      csvColumns,
      ...csvDatas.map((item) => [
        item.Name,
        item.Magnification,
        item.Watt,
        item.Watt * item.Magnification,
        item.LatestUpdate,
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
          {csvDatas.map((device, index) => (
            <tr key={index}>
              <td>
                <input type="checkbox" onChange={handleCheckboxChange} />
              </td>
              <td>{device.Name}</td>
              <td>{device.Magnification}</td>
              <td>{device.Watt}</td>
              <td>{device.Watt * device.Magnification}</td>
              <td>
                {new Date(Date.parse(device.LatestUpdate))
                  .toISOString()
                  .replace("T", " ")
                  .substring(0, 19)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
