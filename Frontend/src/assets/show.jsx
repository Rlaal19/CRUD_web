import { faPenToSquare, faTrash } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from 'react-router-dom';

function Show() {
    const [users, setUsers] = useState([]); // เก็บข้อมูลที่ดึงมาจาก Backend
    const navigate = useNavigate();

    // ดึงข้อมูลจาก Backend
    useEffect(() => {
        fetch("http://localhost:8000/humans")
            .then((response) => response.json())
            .then((data) => setUsers(data))
            .catch((error) => console.error("Error fetching data:", error));
    }, []);

    // ฟังก์ชันลบข้อมูล
    const handleDelete = (id) => {
        const confirmDelete = window.confirm("Are you sure you want to delete this user?");
        if (confirmDelete) {
            fetch(`http://localhost:8000/humans/${id}`, {
                method: "DELETE",
            })
                .then((response) => {
                    if (response.ok) {
                        // ลบข้อมูลใน state
                        setUsers(users.filter((user) => user.id !== id));
                    } else {
                        console.error("Failed to delete user");
                    }
                })
                .catch((error) => console.error("Error deleting user:", error));
        }
    };

    return (
        <div>
            <div className="grid grid-cols-3 mt-8 ml-10 mr-10 py-4 bg-slate-400 opacity-80 rounded-md">
                <label className="ml-3 mr-2 font-medium">Name</label>
                <label className="mr-2 font-medium">Surname</label>
                <label className="mr-2 font-medium">Actions</label>
            </div>
            {users.map((user) => (
                <div key={user.id} className="grid grid-cols-3 ml-12 mr-12 py-2">
                    <div className="ml-2 mt-2">{user.F_name}</div>
                    <div className="mt-2">{user.L_name}</div>
                    <div>
                        <Link className="text-error" to={`/edit/${user.id}`}>
                            <button className="mr-6 btn btn-warning">
                                Edit <FontAwesomeIcon icon={faPenToSquare} />
                            </button>
                        </Link>
                        <button
                            className="btn btn-warning"
                            onClick={() => handleDelete(user.id)}
                        >
                            Delete <FontAwesomeIcon icon={faTrash} />
                        </button>
                    </div>
                </div>
            ))}
        </div>
    );
}

export default Show;
