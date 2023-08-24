import { useMemo, useRef, useState} from 'react';
import {ActionIcon,  Grid, Modal} from "@mantine/core";

import {AgGridReact} from 'ag-grid-react'; // the AG Grid React Component

import 'ag-grid-community/styles/ag-grid.css'; // Core grid CSS, always needed
import 'ag-grid-community/styles/ag-theme-alpine.css'; // Optional theme CSS

import DashBoardHeader from "./components/DashboardHeader.tsx";
import {useQuery} from "@tanstack/react-query";
import {

    IconSettings,

} from '@tabler/icons-react';
import {useDisclosure} from "@mantine/hooks";

interface ManagerDetail {
    serverDetail: { clients: CLientData[], numberOfClients: number };
}

interface CLientData {
    client_id: string,
    login_in_at: string
}

function Dashboard() {
    const [opened, {open, close}] = useDisclosure(false);
    const [selectedClient, setSelectedClient] = useState<CLientData>();

    const gridRef = useRef<any>(); // Optional - for accessing Grid's API
    // Each Column Definition results in one Column.
    const [columnDefs] = useState<any>([
        {field: 'client_id', filter: true, headerName: "Client ID"},
        {field: 'login_in_at', filter: true, width: 300, headerName: "Login In Since"},
        {
            headerName: "Tools", cellRenderer: (parms: any) => {
                console.log(parms.data)
                return (
                    <ActionIcon variant="subtle" className={'mt-1'} onClick={() => {
                        open()
                        setSelectedClient(parms.data)
                    }}><IconSettings size="1rem"/></ActionIcon>
                )
            }
        }

    ]);

    // DefaultColDef sets props common to all Columns
    const defaultColDef = useMemo(() => ({
        sortable: true
    }), []);


    const {isLoading, data} = useQuery({
            queryKey: ['clients'], queryFn: async () => {
                try {
                    const response = await fetch("/admin/stats");
                    if (response.ok) {
                        const data: ManagerDetail = await response.json();
                        console.log("Data fetched:", data);
                        return data.serverDetail.clients
                    } else {
                        console.log("Fetch Error:", response.statusText);
                    }
                } catch (error) {
                    console.error("Error fetching data:", error);
                }
                return []
            },
        }
    )


    return (
        <Grid gutter={5}  gutterXl={0} className={''}>
            <Grid.Col span={12}>
                <Modal opened={opened} onClose={close} title="CLient Connection Control" size={700}>
                    {JSON.stringify(selectedClient)}
                </Modal>
                <DashBoardHeader/>
                <div className="ag-theme-alpine px-10 rounded-md h-80" >
                    <AgGridReact
                        className={'rounded-md'}
                        ref={gridRef} // Ref for accessing Grid's API
                        rowData={isLoading ? [] : data} // Row Data for Rows
                        columnDefs={columnDefs} // Column Defs for Columns
                        defaultColDef={defaultColDef} // Default Column Properties
                        animateRows={true} // Optional - set to 'true' to have rows animate when sorted
                        rowSelection='multiple' // Options - allows click selection of rows
                    />
                </div>
            </Grid.Col>
        </Grid>

    );
}

export default Dashboard;
