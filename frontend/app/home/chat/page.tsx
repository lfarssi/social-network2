"use client";

import { useState } from "react";
import ChatList from "@/components/chatList";
import ChatWindow from "@/components/chatWindow";
import { Chat } from "@/types/chat";

export default function ChatPage() {
  const [activeChat, setActiveChat] = useState<Chat | null>(null);

  return (
    <div className="chat-page">
      <div className="chat-container">
        <ChatList activeChat={activeChat} setActiveChat={setActiveChat} />
        <ChatWindow chat={activeChat} />
      </div>
    </div>
  );
}
