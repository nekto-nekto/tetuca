create or replace function after_mod_log_insert()
returns trigger as $$
declare
	op bigint;
begin
	if new.post_id != 0 then
		insert into post_moderation (post_id, type, "by", length, data)
			values (new.post_id, new.type, new."by", new.length, new.data);
		update posts
			set moderated = true
			where id = new.post_id
			returning posts.op into op;
		perform pg_notify('post_moderated',
			concat_ws(',', op, new.id));

		-- Posts bump threads only on creation and closure
		if (new.type = 2) then
			perform bump_thread(op, true, true, new.post_id);
		--else
		--	perform bump_thread(op, true);
		end if;
	end if;
	return null;
end;
$$ language plpgsql;
