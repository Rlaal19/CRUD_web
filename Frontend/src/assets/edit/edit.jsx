import React, { useState, useEffect } from 'react';
import './edit.css';
import { useParams, useNavigate } from 'react-router-dom';

const Edit = () => {
    const { id } = useParams(); // รับ id จาก URL
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        F_name: '',
        L_name: '',
    });

    // ดึงข้อมูลจาก Backend เพื่อแสดงในฟอร์ม
    useEffect(() => {
        fetch(`http://127.0.0.1:8000/humans/${id}`)
            .then((response) => response.json())
            .then((data) => {
                setFormData({
                    F_name: data.F_name,
                    L_name: data.L_name,
                });
            })
            .catch((error) => console.error("Error fetching data:", error));
    }, [id]);

    // อัปเดต state เมื่อมีการเปลี่ยนแปลงในฟอร์ม
    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: value,
        });
    };

    // ส่งคำขอแก้ไขข้อมูลไปยัง Backend
    const handleSubmit = (e) => {
        e.preventDefault();
        fetch(`http://nginx:8000/humans/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData),
        })
            .then((response) => {
                if (response.ok) {
                    navigate('/'); // กลับไปหน้าแรกหลังแก้ไขสำเร็จ
                } else {
                    console.error("Failed to update user");
                }
            })
            .catch((error) => console.error("Error updating user:", error));
    };

    return (
        <div className='wrapper'>
            <form onSubmit={handleSubmit}>
                <div className='pl-5 pr-5 pb-10 bg-slate-200 rounded-lg opacity-90'>
                    <h1>Edit</h1>
                    <div className='inputbox'>
                        <input
                            type="text"
                            name="F_name"
                            value={formData.F_name}
                            placeholder="Name"
                            className="input input-bordered input-warning w-full max-w-xs"
                            onChange={handleChange}
                        />
                    </div>

                    <div className='inputbox'>
                        <input
                            type="text"
                            name="L_name"
                            value={formData.L_name}
                            placeholder="Surname"
                            className="input input-bordered input-warning w-full max-w-xs"
                            onChange={handleChange}
                        />
                    </div>

                    <div className='flex justify-center signup-btn'>
                        <button type="submit" className="btn btn-success">Change</button>
                    </div>
                </div>
            </form>
        </div>
    );
};

export default Edit;
